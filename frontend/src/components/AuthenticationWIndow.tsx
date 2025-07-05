import "./AuthenticationWindow.css"
import { useState } from "react"

function AuthenticationWindow() {
  const [inputbox, setInputbox] = useState<string>("")

  return <>
    <div className="container">
      <label className="enterPassword">Please enter your password</label>
      <div className="lineBreak"/>

      <input className="bigInput" 
        onChange={(e) =>{setInputbox(e.target.value)}}
        onKeyDown={(e) => {
          if(e.code == "Enter") {
            // @ts-ignore
            localStorage.setItem("password", inputbox)
            window.location.reload()
          }
        }}/>

      <button className='loginButton' onClick={() => {
        localStorage.setItem("password", inputbox)
        window.location.reload()
      }}>Log In</button>

    </div>
  </>
}

export default AuthenticationWindow