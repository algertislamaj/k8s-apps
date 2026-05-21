import { useState } from "react"

export default function App() {
  const [message, setMessage] = useState("")

  const handleClick = async () => {
    const res = await fetch("http://localhost:8080/api/hello")
    const data = await res.json()
    setMessage(data.message)
  }

  return (
    <div style={{ padding: "2rem", fontFamily: "sans-serif" }}>
      <h1>VM Manager</h1>
      <button onClick={handleClick}>Say Hello</button>
      {message && <p>{message}</p>}
    </div>
  )
}
