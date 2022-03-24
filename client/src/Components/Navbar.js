import {BrowserRouter,Route,Switch,Link} from "react-router-dom";
function Navbar(){
    return (<div className="sticky top-0">
            <nav>
                <ul className="bg-[#4c1d95] flex justify-around h-10 items-center text-[#67e8f9]">
                    <li className="inline-block">
                        <Link to="/dashboard">Dashboard</Link>
                    </li>
                    <li className="inline-block">
                        <Link to="/transactions">Payments</Link>
                    </li>
                    <li className="inline-block">
                        <Link to="/expenses">Expenses</Link>
                    </li>
                    <li className="inline-block">
                        <Link to="/logout">Logout</Link>
                    </li>
                </ul>
            </nav>
        </div>
    )  
}
export default Navbar;