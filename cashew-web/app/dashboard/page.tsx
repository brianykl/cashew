
import { DashboardContent } from "../components/dashboardContent"
import { getAccessToken, getSession } from "@auth0/nextjs-auth0"

type SearchParams = { [key: string]: string | string[] | undefined };

type DashboardProps = {
    searchParams: SearchParams;
};

export default async function Dashboard({ searchParams }: DashboardProps) {
    const publicToken = searchParams.public_token
    let exchangeResult = null   

    if (publicToken) {
        try {
            const { accessToken } = await getAccessToken();
            const response = await fetch('http://localhost:8080/protected/exchange', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify({public_token: publicToken})
            });
            exchangeResult = await response.text()
            console.log('Token exchanged successfully:', exchangeResult);
        } catch (error) {
            console.error('Error exchanging token:', error);
        }
    }
  
    return (
        <DashboardContent />
    )
}