import "./Ribbon.css"
import type { JSX } from "react";
import type { UserInfo } from "../models/User";

function Ribbon({userInfo}: {userInfo: UserInfo | null}): JSX.Element {
    return (<>
        <div className="Ribbon">
            Server Manager
        </div>
    </>);
}

export default Ribbon