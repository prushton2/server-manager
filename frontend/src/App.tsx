import './App.css'
import { useEffect, useState, type JSX } from 'react'
import type { ServerState, State } from './models/State'
import { Action, Authenticate, GetState } from './API'
import { toast } from "react-fox-toast"
import type { UserInfo } from './models/User'
import Ribbon from './components/Ribbon'
import AuthenticationWindow from './components/AuthenticationWIndow'

function App() {
  const [state, setState] = useState<State | null>(null)
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null)
  const [requiresAuth, setRequiresAuth] = useState<boolean>(false)

  useEffect(() => {
    async function init() {
      let password: string = localStorage.getItem("password")+""
      let auth = await Authenticate(password)

      if(typeof auth === "object") {
        setUserInfo(auth)
      } else {
        toast.error(`Please authenticate before using the site.`)
        setRequiresAuth(true)
      }

      let newstate = await GetState(password)
      setState(newstate)
    }
    init()
  }, [])

  function renderState() {
    if (state == null) {
      return <>state is null</>
    }
    if (userInfo == null) {
      return <>userInfo is null</>
    }

    return Object.entries(state.servers).map(([name, serverState]) => {
        return <Server key={name} name={name} serverState={serverState} userInfo={userInfo}/>
      })
  }
  
  return (
    <>
      <Ribbon userInfo={userInfo} />
      {requiresAuth ? <AuthenticationWindow /> : renderState()}
    </>
  )
}

export default App

function Server({name, serverState, userInfo}: {name: string, serverState: ServerState, userInfo: UserInfo}): JSX.Element {
  const [timeRemaining, setTimeRemaining] = useState<number>(0)

  useEffect(() => {
    setTimeRemaining(serverState.endsAt - new Date().getTime() / 1000)
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
      { timeRemaining > 0 && userInfo.canStop ?
        <button className="stopButton" onClick={() => serverAction(name, "stop")}>Stop</button> : <></>
      }
      { timeRemaining > 0  && userInfo.canExtend ?
        <button className="serverButton" onClick={() => serverAction(name, "extend")}>Extend</button> : <></>
      }
      { timeRemaining <= 0 && userInfo.canStart ?
        <button className="serverButton" onClick={() => serverAction(name, "start")}>Start</button> : <></>
      }
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

async function serverAction(name: string, action: string) {
  let password: string = localStorage.getItem("password")+""

  let response = await Action(name, password, action)

  if(response == "") {
    window.location.reload()
  } else {
    toast.error(`Error: ${response}`)
  }
}