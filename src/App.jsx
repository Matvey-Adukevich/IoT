import { useState } from 'react'
import './App.css'

const BASE = 'http://localhost:8080'

export default function App() {
  const [temp,     setTemp]     = useState(null)
  const [windows,  setWindows]  = useState(null)
  const [computer, setComputer] = useState(null)
  const [logs,     setLogs]     = useState(null)
  const [loading,  setLoading]  = useState(false)

  async function load(url, setValue) {
    setLoading(true)
    try {
      const res = await fetch(BASE + url)
      setValue(await res.json())
    } catch {
      setValue('Ошибка')
    }
    setLoading(false)
  }

  async function loadAll() {
    await load('/temperature', setTemp)
    await load('/windows',     setWindows)
    await load('/pcstatus',    setComputer)
  }

  return (
    <div className="dashboard">

      <div className="dashboard-row">

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Температура</span>
            <span className="panel-value">{temp ?? '—'}</span>
          </div>
          <button type="button" className="panel-btn" disabled={loading}
            onClick={() => load('/temperature', setTemp)}>
            Узнать температуру
          </button>
        </div>

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Статус окон</span>
            <span className="panel-value">{windows ?? '—'}</span>
          </div>
          <button type="button" className="panel-btn" disabled={loading}
            onClick={() => load('/windows', setWindows)}>
            Узнать статус
          </button>
        </div>

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Статус компьютеров</span>
            <div className="panel-value">
              {computer
                ? Object.entries(computer).map(([id, status]) => (
                    <div key={id} className={`pc-status pc-status--${status}`}>
                      ПК №{id}: {status.toUpperCase()}
                    </div>
                  ))
                : '—'}
            </div>
          </div>
          <button type="button" className="panel-btn" disabled={loading}
            onClick={() => load('/pcstatus', setComputer)}>
            Узнать статус
          </button>
        </div>

      </div>

      <div className="dashboard-row">

        <div className="panel panel--large">
          <div className="panel-display">
            <span className="panel-title">Обновить всё сразу</span>
            {loading && <span className="panel-loading">Загрузка...</span>}
          </div>
          <button type="button" className="panel-btn" disabled={loading} onClick={loadAll}>
            Получить статус всего
          </button>
        </div>

        <div className="panel panel--logs">
          <div className="panel-display">
            <span className="panel-title">Логи</span>
            <div className="logs-box">
              {logs?.length > 0
                ? logs.map((log, i) => <div key={i} className="log-line">{log}</div>)
                : '—'}
            </div>
          </div>
          <button type="button" className="panel-btn" disabled={loading}
            onClick={() => load('/logs', setLogs)}>
            Получить логи
          </button>
        </div>

      </div>

    </div>
  )
}
