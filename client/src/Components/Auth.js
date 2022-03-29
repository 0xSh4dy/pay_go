import React from 'react';
import ReactDOM from 'react-dom';
import {useState} from 'react';
import { Navigate, useNavigate } from 'react-router-dom';
import { useEffect } from 'react';
import axios from 'axios';
import Cookies from 'universal-cookie';
import Dashboard from './Dashboard';

function ResetMessage(setCustomMessage,setMsgBoxColor,setMsgTxtColor){
    setTimeout(()=>{
        setCustomMessage("");
        setMsgBoxColor("");
        setMsgTxtColor("");
    },5000);
}
function Home(){
    const [mode,setMode] = useState("login");
    const [modeText,setModeText] = useState("Do not have an account? Register");
    const [username,setUsername] = useState("");
    const [password,setPassword] = useState("");
    const [email,setEmail] = useState("");
    const [loginOK,setLoginOK] = useState(false);
    const [signupOK,setSignupOK] = useState(false);
    const [customMessage,setCustomMessage] = useState("");
    const [msgBoxColor,setMsgBoxColor] = useState("");
    const [msgTxtColor,setMsgTxtColor] = useState("");
    useEffect(()=>{
        const cookie = new Cookies();
        let authCookie = cookie.get("auth");
        if(authCookie!==undefined){
            setLoginOK(true);
        }
        else{
            setLoginOK(false);
        }
    },[])
    function handleLogin(){
        let data = {
            username:username,
            password:password
        }
        axios({
            method:"post",
            url:"http://127.0.0.1:7000/api/login/",
            data:data
        }).then((res)=>{
            console.log(res.data);
            if(res.data=="Invalid username or password"){
                setCustomMessage("Invalid credentials");
                setMsgBoxColor("red");
                setMsgTxtColor("#F5F106");
                ResetMessage(setCustomMessage,setMsgBoxColor,setMsgTxtColor);
            }
            else if(res.data=="Internal server error"){
                setCustomMessage("Internal server error");
                setMsgBoxColor("#F5E306");
                setMsgTxtColor("#23F506");
                ResetMessage(setCustomMessage,setMsgBoxColor,setMsgTxtColor);
            }
            else{
                const cookies = new Cookies();
                cookies.set("auth",res.data,"/");
                window.location.href = "/dashboard";
                setLoginOK(true);
            }
        });
    }
    function handleRegister(){
        let data = {
            username:username,
            password:password,
            email:email
        }
        axios({
            method:"post",
            url:"http://127.0.0.1:7000/api/register/",
            data:data
        }).then(
           (res)=>{
               console.log(res.data)
               if(res.data==="Email is already taken"){
                    setCustomMessage(res.data);
                    setMsgBoxColor("red");
                    setMsgTxtColor("#F5F106");
                    ResetMessage(setCustomMessage,setMsgBoxColor,setMsgTxtColor);
               }
               else if(res.data==="Username is already taken"){
                    setCustomMessage(res.data);
                    setMsgBoxColor("red");
                    setMsgTxtColor("#F5F106");
                    ResetMessage(setCustomMessage,setMsgBoxColor,setMsgTxtColor);
               }
               else if(res.data==="Error"){
                    setCustomMessage("Internal server error");
                    setMsgBoxColor("#F5E306");
                    setMsgTxtColor("#23F506");
                    ResetMessage(setCustomMessage,setMsgBoxColor,setMsgTxtColor);
               }
               else{
                    setCustomMessage("Successfully registered!!");
                    setMsgBoxColor("#43F506");
                    setMsgTxtColor("#068CF5 ");
                    ResetMessage(setCustomMessage,setMsgBoxColor,setMsgTxtColor);
                    setSignupOK(true);
               }
           }
        );
    }
    function Login(){
        return (<div className='loginBox text-center'>
            <div className='msgBox w-fit p-2' style={{backgroundColor:msgBoxColor,color:msgTxtColor}}>{customMessage}</div>
        <div className="p-5 loginForm  text-blue-600">
            <p className="text-2xl">Login</p>
            <input type="text" placeholder="Enter username" id="loginName" className="inp block p-2 mt-2" onChange={(e)=>setUsername(e.target.value)}/>
            <input type="password" placeholder="Enter password" id="loginPassword" className="inp block mt-5 p-2 " onChange={(e)=>setPassword(e.target.value)}/>
            <button onClick={handleLogin} id="loginBtn" className="p-2 mt-5 bg-[#a21caf] block rounded-md text-[#EDFC0F]" >Login!</button>
            <p className='mt-5 txtHover' onClick={()=>{
                setMode("register");
                setModeText("Already have an account? Login");
            }}>{modeText}</p>
            <p className='mt-5 txtHover' onClick={()=>{
                window.location.href = "/change"
            }}>Forgot password? Click here to reset</p>
        </div>  </div>)  
    }
    function Signup(){
        return (<div className='signupBox'>
            <div className='msgBox w-fit p-2' style={{backgroundColor:msgBoxColor,color:msgTxtColor}}>{customMessage}</div>
        <div className="p-5 signupForm text-blue-600">
            <p className="text-2xl">Register</p>
            <input type="text" placeholder="Enter username" id="signupName" className="inp block p-2 mt-2" onChange={(e)=>{setUsername(e.target.value)}}/>
            <input type="password" placeholder="Enter password" id="signupPassword" className="inp block mt-5 p-2 " onChange={(e)=>{setPassword(e.target.value)}}/>
            <input type="email" placeholder="Enter email" id='signupEmail' className='inp block p-2 mt-5' onChange={(e)=>{setEmail(e.target.value)}}/>
            <button onClick={handleRegister} id="signupBtn" className="p-2 mt-5 bg-[#a21caf] block rounded-md text-[#EDFC0F]" >Register!</button>
            <p className='mt-5 txtHover' onClick={
                ()=>{
                    setMode("login");
                setModeText("Do not have an account? Register");
                }
            } >{modeText}</p>
        </div>  
        </div>)
    }
    if(loginOK==true){
        window.location.href = "/dashboard";
    }
    else if(signupOK===true){
        return <div>
            {Login()}
        </div>
    }
    return <div>
        {mode==="login"?Login():Signup()}
    </div>
}

export default Home;