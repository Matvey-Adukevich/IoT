import { useState } from 'react'
import './App.css'

export default function App() {

  const [temp, setTemp] = useState(null)
  const [windows, setWindows]  = useState(null)
  const [computer, setComputer] = useState(null)
  const [logs, setLogs] = useState(null)

  const [loading, setLoading] = useState(false)

  async function load(url, setValue) {
    setLoading(true)
    try {
      const res = await fetch(url)
      const value = await res.json()
      setValue(value)
    } catch {
      setValue('Ошибка')
    }
    setLoading(false)
  }

  async function loadAll() {
    await load('/api/status', setTemp)
    await load('/api/windows', setWindows)
    await load('/api/computer', setComputer)
  }

  return (
    <div className="dashboard">

      <div className="dashboard-row">

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Температура</span>
            <span className="panel-value">{temp}</span>
          </div>
          <button className="panel-btn" disabled={loading} onClick={() => load('/api/temperature', setTemp)}>Узнать температуру</button>
        </div>

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Статус окон</span>
            <span className="panel-value">{windows}</span>
          </div>
          <button className="panel-btn" disabled={loading} onClick={() => load('/api/windows', setWindows)}>Узнать статус</button>
        </div>

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Статус компьютеров</span>
            <span className="panel-value">{computer}</span>
          </div>
          <button className="panel-btn" disabled={loading} onClick={() => load('/api/computer', setComputer)}> Узнать статус </button>
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

        <div className="panel">
          <div className="panel-display">
            <span className="panel-title">Логи</span>
            <span className="panel-value">{logs}</span>
          </div>
          <button className="panel-btn" disabled={loading} onClick={() => load('/api/logs', setLogs)}>Получить логи</button>
        </div>

      </div>

    </div>
  )
}
