import { useState } from 'react'
import './App.css'

const BASE_URL = 'http://localhost:8080'

export default function App() {
  const [temp, setTemp] = useState(null)
  const [windows, setWindows] = useState(null)
  const [computer, setComputer] = useState(null)
  const [logs, setLogs] = useState(null)

  const [loading, setLoading] = useState(false)

  async function load(url, parserFn) {
    setLoading(true)
    try {
      const res = await fetch(`${BASE_URL}${url}`)
      if (!res.ok) throw new Error()
      const value = await res.json()
      parserFn(value)
    } catch {
      alert(`Не удалось загрузить данные с ${url}`)
    }
    setLoading(false)
  }

  async function loadAll() {
    setLoading(true)
    try {
      const res = await fetch(`${BASE_URL}/status`)
      if (!res.ok) throw new Error()
      const data = await res.json()
      
      setTemp(data.temperature)
      setWindows(data.windows)
      setComputer(data.pcs)
    } catch {
      alert('Ошибка при генерации общего статуса')
    }
    setLoading(false)
  }

  return (
    <div className="dashboard">

      <div className="dashboard-row">

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Температура</span>
            <span className="panel-value">
              {temp !== null ? (typeof temp === 'object' ? `${temp.temperature} °C` : `${temp} °C`) : '—'}
            </span>
          </div>
          <button className="panel-btn" disabled={loading} onClick={() => load('/temperature', setTemp)}>Узнать температуру</button>
        </div>

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Статус окон</span>
            <span className="panel-value">
              {windows !== null ? (typeof windows === 'object' ? windows.windows : windows) : '—'}
            </span>
          </div>
          <button className="panel-btn" disabled={loading} onClick={() => load('/windows', setWindows)}>Узнать статус</button>
        </div>

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Статус компьютеров</span>
            <div className="panel-value" style={{ fontSize: '16px', maxHeight: '100px', overflowY: 'auto' }}>
              {computer ? (
                Object.entries(computer).map(([id, status]) => (
                  <div key={id} style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                    <span>ПК №{id}:</span>
                    <span style={{ color: status === 'on' ? '#04d361' : '#f75a68', fontWeight: 'bold' }}>{status.toUpperCase()}</span>
                  </div>
                ))
              ) : '—'}
            </div>
          </div>
          <button className="panel-btn" disabled={loading} onClick={() => load('/pcstatus', setComputer)}> Узнать статус </button>
        </div>

      </div>


      <div className="dashboard-row">
        <div className="panel panel--large">
          <div className="panel-display">
            <span className="panel-title">Обновить всё сразу</span> 
            {loading && <span className="panel-loading">Загрузка...</span>}
          </div>
          <button className="panel-btn" disabled={loading} onClick={loadAll}>Получить статус всего</button>
        </div>
        <div className="panel panel--logs">
          <div className="panel-display">
            <span className="panel-title">Логи</span>
            <div className="logs-box">
              {logs?.length > 0
                ? logs.map((log, i) => <div key={i} className="logs-line">{log}</div>)
                : <span style={{ color: 'inherit' }}>—</span>}
            </div>
          </div>
          <button className="panel-btn" disabled={loading}
            onClick={() => load('/logs', setLogs)}>
            Получить логи
          </button>
        </div>

      </div>

    </div>
  )
}
