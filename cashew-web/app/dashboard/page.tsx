import { redirect } from "next/navigation"
import { DashboardContent } from "../components/dashboardContent"
import { getAccessToken, getSession } from "@auth0/nextjs-auth0"


export default async function Dashboard({searchParams,}: {searchParams: { [key: string]: string | string[] | undefined}}) {
    const session = await getSession()
    const accessToken = await getAccessToken()
    if (!session) {
        redirect('/api/auth/login');
    }

    let exchangeResult = null;
    const publicToken = searchParams.public_token;
    if (typeof publicToken === 'string') {
        try {
          exchangeResult = await fetch('http://localhost:8080/exchange', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${accessToken}` 
            },
            body: JSON.stringify({public_token: publicToken})
          });
          console.log('Token exchanged successfully:', exchangeResult.json());
        } catch (error) {
            console.error('Error exchanging token:', error);
        }
    }      
    return (
        <DashboardContent />
    )
}