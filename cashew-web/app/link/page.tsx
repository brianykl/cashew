import { getSession, getAccessToken } from "@auth0/nextjs-auth0"
import { redirect } from "next/navigation"
import { LinkContent } from "../components/linkContent"

export default async function Link() {
    const session = await getSession()
    const { accessToken } = await getAccessToken()

    if(!session || !accessToken) {
        redirect('api/auth/login')
    }

    const response = await fetch('http://localhost:8080/link')
    const data = await response.json()
    const linkToken = data.link_token

    return (
        <LinkContent initialLinkToken={linkToken} accessToken={accessToken!}/>
    )
}