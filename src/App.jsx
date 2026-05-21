import { useState } from 'react'
import './App.css'

const BASE_URL = 'http://localhost:8080'

async function fetchJSON(url) {
  const res = await fetch(BASE_URL + url)
  if (!res.ok) throw new Error()
  return res.json()
}

export default function App() {
  const [temp,     setTemp]     = useState(null)
  const [windows,  setWindows]  = useState(null)
  const [computer, setComputer] = useState(null)
  const [logs,     setLogs]     = useState(null)
  const [loading,  setLoading]  = useState(false)

  async function load(url, setter) {
    setLoading(true)
    try { setter(await fetchJSON(url)) }
    catch { alert(`Не удалось загрузить данные с ${url}`) }
    finally { setLoading(false) }
  }

  async function loadAll() {
    setLoading(true)
    try {
      const data = await fetchJSON('/status')
      setTemp(data.temperature)
      setWindows(data.windows)
      setComputer(data.pcs)
    } catch {
      alert('Ошибка при получении общего статуса')
    } finally {
      setLoading(false)
    }
  }

  const tempValue = temp === null ? '—'
    : `${typeof temp === 'object' ? temp.temperature : temp} °C`

  const windowsValue = windows === null ? '—'
    : typeof windows === 'object' ? windows.windows : windows

  return (
    <div className="dashboard">
      <div className="dashboard-row">

        <Panel title="Температура" onLoad={() => load('/temperature', setTemp)} loading={loading}>
          <span className="panel-value">{tempValue}</span>
        </Panel>

        <Panel title="Статус окон" onLoad={() => load('/windows', setWindows)} loading={loading}>
          <span className="panel-value">{windowsValue}</span>
        </Panel>

        <Panel title="Статус компьютеров" onLoad={() => load('/pcstatus', setComputer)} loading={loading} btnText="Узнать статус">
          <div className="panel-value" style={{ fontSize: 16, maxHeight: 100, overflowY: 'auto' }}>
            {computer
              ? Object.entries(computer).map(([id, status]) => (
                  <div key={id} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                    <span>ПК №{id}:</span>
                    <span style={{ color: status === 'on' ? '#04d361' : '#f75a68', fontWeight: 'bold' }}>
                      {status.toUpperCase()}
                    </span>
                  </div>
                ))
              : '—'}
          </div>
        </Panel>

      </div>
      <div className="dashboard-row">

        <Panel title="Обновить всё сразу" onLoad={loadAll} loading={loading} btnText="Получить статус всего">
          {loading && <span className="panel-loading">Загрузка...</span>}
        </Panel>

        <Panel title="Логи" onLoad={() => load('/logs', setLogs)} loading={loading} btnText="Получить логи"
          style={{ width: 504 }}>
          <div style={{ background: '#000', color: '#aaa', fontFamily: 'monospace', fontSize: 11,
            padding: 8, borderRadius: 4, maxHeight: 120, overflowY: 'auto' }}>
            {logs?.length > 0
              ? logs.map((log, i) => (
                  <div key={i} style={{ whiteSpace: 'nowrap', borderBottom: '1px solid #222', paddingBottom: 4, textAlign: 'left' }}>
                    {log}
                  </div>
                ))
              : <span style={{ color: 'inherit' }}>—</span>}
          </div>
        </Panel>

      </div>
    </div>
  )
}

function Panel({ title, onLoad, loading, btnText = 'Узнать статус', style, children }) {
  return (
    <div className="panel" style={style}>
      <div className="panel-display">
        <span className="panel-title">{title}</span>
        {children}
      </div>
      <button className="panel-btn" disabled={loading} onClick={onLoad}>{btnText}</button>
    </div>
  )
}