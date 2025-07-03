import './App.css'
import { useEffect, useState, type JSX } from 'react'
import type { ServerState, State } from './models/State'
import { Authenticate, Extend, GetState, Start } from './API'
import { toast } from "react-fox-toast"

function App() {
  const [state, setState] = useState<State | null>(null)
  const [requiresAuth, setRequiresAuth] = useState<boolean>(false)

  useEffect(() => {
    async function init() {
      let password: string = localStorage.getItem("password")+""
      let auth = await Authenticate(password)

      if(auth != "") {
        toast.error(`Please authenticate before using the site.`)
        setRequiresAuth(true)
      }

      let newstate = await GetState(password)
      setState(newstate)
    }
    init()
  }, [])

  function renderState() {
    return <>
      {
        state == null ? <>State is null</> : 
        Object.entries(state.servers)
        .map(([name, serverState]) => {
          return <Server key={name} name={name} serverState={serverState} />
        })
      }
    </>
  }

  return (
    <>
      {requiresAuth ? AuthenticationWindow() : renderState()}
    </>
  )
}

export default App

function AuthenticationWindow() {
  return <>
    <div className="serverDiv">
      <label className="enterPassword">Please enter your password</label>
      <div className="lineBreak"/>
      <input className="bigInput" onKeyDown={(e) => {
          if(e.code == "Enter") {
            // @ts-ignore
            localStorage.setItem("password", e.target.value)
            window.location.reload()
          }
        }}/>
    </div>
  </>
}

function Server({name, serverState}: {name: string, serverState: ServerState}): JSX.Element {
  const [timeRemaining, setTimeRemaining] = useState<number>(0)

  useEffect(() => {
    const interval = setInterval(() => {
      setTimeRemaining(serverState.endsAt - new Date().getTime() / 1000)
    }, 1000)
    return () => clearInterval(interval)
  }, [])
  // @ts-ignore is it really this dumb? its not my fault theres no real int type

  return <div className="serverDiv">
    <label className="serverName">{name.charAt(0).toUpperCase() + name.slice(1)}</label>
    <label className="serverStatus"> {timeRemaining < 0 ? "Server is off" : `${formatTime(timeRemaining)}`}</label>
    <div className="serverButtonContainer">
      <button className="serverButton" onClick={() => startOrExtend(name, timeRemaining)}>{timeRemaining < 0 ? "Start" : "Extend"}</button>
    </div>
  </div>

}

function formatTime(time: number): string {
  // @ts-ignore
  let seconds = parseInt(time % 60)
  // @ts-ignore
  let minutes = parseInt((time / 60)%60)
  // @ts-ignore
  let hours =   parseInt((time / 60 / 60)%60)
  // @ts-ignore
  let days =    parseInt((time / 60 / 60 / 24))

  let secondsStr = String(seconds).padStart(2, '0')
  let minutesStr = String(minutes).padStart(2, '0')
  let hoursStr = String(hours).padStart(2, '0')
  let daysStr = String(days).padStart(2, '0')

  return `${(daysStr)}:${(hoursStr)}:${(minutesStr)}:${(secondsStr)}`
}

async function startOrExtend(name: string, timeRemaining: number) {
  let password: string = localStorage.getItem("password")+""
  if(timeRemaining < 0) {

    let response = await Start(name, password)
    if(response == "") {
      toast.success('Starting Server')
      window.location.reload()
    } else {
      toast.error(`Error: ${response}`)
    }

  } else {

    let response = await Extend(name, password)
    if(response == "") {
      toast.success('Extending Server')
      window.location.reload()
    } else {
      toast.error(`Error: ${response}`)
    }
  }
}