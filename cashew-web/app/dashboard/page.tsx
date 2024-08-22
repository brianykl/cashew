import { getSession } from "@auth0/nextjs-auth0";
import { redirect } from "next/navigation";
import { DashboardContent } from "../components/dashboardContent";

export default async function Dashboard() {
    const session = await getSession();

    if(!session) {
        redirect('api/auth/login')
    }

    const response = await fetch('http://localhost:8080/link');
    const data = await response.json();
    const linkToken = data.link_token

    return (
        <DashboardContent initialLinkToken={linkToken} />
    );
}