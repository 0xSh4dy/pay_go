import {BrowserRouter,Routes,Route} from "react-router-dom";
import Dashboard from "./Dashboard";
import Auth from "./Auth";
import Transcations from "./Transcations";
import Expenses from "./Expenses";
import Logout from "./Logout";

function Router(){
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Auth/>}/>
                <Route path="dashboard" element={<Dashboard/>}/>
                <Route path="transactions" element={<Transcations/>}/>
                <Route path="expenses" element={<Expenses/>}/>
                <Route path="logout" element={<Logout/>}/>
            </Routes>
        </BrowserRouter>
    )
}
export default Router;