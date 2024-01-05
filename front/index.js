import React, { useEffect, useState } from "react"

export default function Home() {
  const [socketInit, setSocketInit] = useState(false)
  const [spreads, setSpreads] = useState([])

  const socketInitializer = () => {
    console.log("init socket")
    let socket = new WebSocket("ws://localhost:4000/bestspreads")
    socket.onmessage = (event) => {
      console.log(JSON.parse(event.data)) // <<<<<<<<<<<<<<<<
      setSpreads(JSON.parse(event.data))
    }
    setSocketInit(true)
  }

  const show = (data) => {
    return Object.keys(data).map((ticker) => {
      return (
        <div id={ticker} className="mb-6">
          <h1 className="text-2xl mb-2 font-bold border-b border-b-1 py-2">
            {ticker}
          </h1>
          <div className="flex justify-between">
            {data[ticker].map((d) => {
              const color = d.Spread > 4 ? "bg-green-400" : ""
              return (
                <table className={`table-auto ${color}`}>
                  <thead>
                    <tr>
                      <TH>{d.A}</TH>
                      <TH>{d.B}</TH>
                      <TH>spread</TH>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <Td className="text-green-500">
                        {Number(d.BestBid).toFixed(2)}
                      </Td>
                      <Td className="text-red-500">
                        {Number(d.BestAsk).toFixed(2)}
                      </Td>
                      <Td>{d.Spread}</Td>
                    </tr>
                  </tbody>
                </table>
              )
            })}
          </div>
        </div>
      )
    })
  }

  useEffect(() => {
    if (!socketInit) {
      socketInitializer()
    }
  }, [])

  return <div className="container">{show(spreads)}</div>
}

const Th = ({ children, className }) => {
  return <th className={`pr-4 ${className}`}>{children}</th>
}

const Td = ({ children, className }) => {
  return <td className={`text-left pr-4 ${className}`}>{children}</td>
} 