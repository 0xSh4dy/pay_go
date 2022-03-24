import {BrowserRouter,Routes,Route} from "react-router-dom";
import Dashboard from "./Dashboard";
import Auth from "./Auth";
import Transcations from "./Transcations";
import Expenses from "./Expenses";
import Logout from "./Logout";
import Unauthorized from "./Unauthorized";
import Badrequest from "./Badrequest";

function Router(){
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Auth/>}/>
                <Route path="dashboard" element={<Dashboard/>}/>
                <Route path="transactions" element={<Transcations/>}/>
                <Route path="expenses" element={<Expenses/>}/>
                <Route path="logout" element={<Logout/>}/>
                <Route path="badrequest" element={<Badrequest/>}/>
                <Route path="unauthorized" element={<Unauthorized/>}/>
            </Routes>
        </BrowserRouter>
    )
}
export default Router;