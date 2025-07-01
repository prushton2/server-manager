import './App.css'
import { useEffect, useState, type JSX } from 'react'
import type { ServerState, State } from './models/State'
import { getState } from './API'

function App() {
  const [state, setState] = useState<State | null>(null)

  useEffect(() => {
    async function init() {
      let newstate = await getState()
      setState(newstate)
      console.log(Object.keys(newstate.servers))
    }
    init()
  }, [])

  return (
    <>
      {state == null ? <>State is null</> : 
        Object.entries(state.servers)
        .map(([name, serverState]) => {
          return <Server key={name} name={name} serverState={serverState} />
        })
      }
    </>
  )
}

export default App


function Server({name, serverState}: {name: string, serverState: ServerState}): JSX.Element {
  
  let timeRemaining = serverState.endsAt - new Date().getTime() / 1000
  // @ts-ignore is it really this dumb? its not my fault theres no real int type
  timeRemaining = parseInt(timeRemaining / 60 / 60)

  return <div className="serverDiv">
    <label className="serverLabel">{name}</label>
    <label className="serverStatus"> {timeRemaining < 0 ? "Server is off" : `${timeRemaining} hours left`}</label>
    <button className="serverButton"> {timeRemaining < 0 ? "Start" : `Extend`}</button>
  </div>

}