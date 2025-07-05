import "./Ribbon.css"
import type { JSX } from "react";
import type { UserInfo } from "../models/User";

function Ribbon({userInfo}: {userInfo: UserInfo | null}): JSX.Element {
    return (<>
        <div className="Ribbon">
            Server Manager
        </div>
        {userInfo == null ? <></> :
        <button className="User" onClick={() => {
            if(confirm("Are you sure you want to log out?")) {
                localStorage.setItem("password", "")
                window.location.reload()
            }
        }}>
            Hi, {userInfo.name} <br /> 
            <label className="LogoutLabel">Click to logout</label>
        </button>
        }
    </>);
}

export default Ribbon