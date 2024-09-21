import { getSession, getAccessToken } from "@auth0/nextjs-auth0"
import { redirect } from "next/navigation"
import { DashboardContent } from "../components/dashboardContent"

export default async function Dashboard() {
    const session = await getSession()
    const { accessToken } = await getAccessToken()

    if(!session || !accessToken) {
        redirect('api/auth/login')
    }
    const userId = session!.user.sub;
    const response = await fetch('http://localhost:8080/link')
    const data = await response.json()
    const linkToken = data.link_token

    return (
        <DashboardContent initialLinkToken={linkToken} accessToken={accessToken!} userId={userId}/>
    )
}