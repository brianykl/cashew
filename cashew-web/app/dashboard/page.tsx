import { redirect } from "next/navigation"
import { DashboardContent } from "../components/dashboardContent"
import { getAccessToken, getSession } from "@auth0/nextjs-auth0"

type SearchParams = { [key: string]: string | string[] | undefined };

type DashboardProps = {
    searchParams: SearchParams;
};

export default async function Dashboard({ searchParams }: DashboardProps) {
    const publicToken = searchParams.public_token   
    const session = await getSession();
    if (!session || !session.user) {
        redirect('api/auth/login')
    }

    const userId = session.user.sub;
      
    if (publicToken) {
        try {
            const { accessToken } = await getAccessToken();
            const response = await fetch('http://localhost:8080/protected/exchange', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify({
                    public_token: publicToken,
                    user_id: userId
                })
            });
            const exchangeStatus = await response.status
            const exchangeBody = await response.body
            if (exchangeStatus == 200) {
                console.log('Token exchanged successfully');
            } else {
                console.log('Unsuccessful exchange:', exchangeBody);
            }
            
        } catch (error) { 
            console.error('Unauthorized to exchange:', error);
        }
    }
  
    return (
        <DashboardContent />
    )
}