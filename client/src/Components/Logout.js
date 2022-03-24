import { useNavigate } from "react-router-dom";
import Cookies from "universal-cookie";
export default function Logout(){
    const cookie = new Cookies();
    if(cookie.get("auth")!==undefined){
        cookie.remove("auth");
        
    }
    window.location.href = "/";
}