import { useEffect } from "react";
import Cookies from "universal-cookie";
import Navbar from "./Navbar";
function Transcations(){
    useEffect(()=>{
        const cookie = new Cookies();
        const authCookie = cookie.get("auth");
        if(authCookie===undefined){
            window.location.href = "";
        }
    },[])

    return (
        <div>
            <Navbar/>
        </div>
    )
}
export default Transcations;