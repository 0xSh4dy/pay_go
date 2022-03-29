import {BrowserRouter,Routes,Route} from "react-router-dom";
import Dashboard from "./Dashboard";
import Auth from "./Auth";
import Transactions from "./Transactions";
import Expenses from "./Expenses";
import Logout from "./Logout";
import Unauthorized from "./Unauthorized";
import Badrequest from "./Badrequest";
import Forgot from "./Forgot";

function Router(){
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Auth/>}/>
                <Route path="dashboard" element={<Dashboard/>}/>
                <Route path="transactions" element={<Transactions/>}/>
                <Route path="expenses" element={<Expenses/>}/>
                <Route path="logout" element={<Logout/>}/>
                <Route path="badrequest" element={<Badrequest/>}/>
                <Route path="unauthorized" element={<Unauthorized/>}/>
                <Route path="change" element={<Forgot/>}/>
            </Routes>
        </BrowserRouter>
    )
}
export default Router;