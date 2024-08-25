'use client'

import { useEffect, useState } from "react"
import { useSearchParams } from "next/navigation"
import { DashboardContent } from "../components/dashboardContent"
import { getAccessToken, getSession } from "@auth0/nextjs-auth0"


export default async function Dashboard() {
    const [exchangeResult, setExchangeResult] = useState(null);
    const searchParams = useSearchParams();
    const publicToken = searchParams.get('public_token');

    useEffect(() => {
        async function exchangeToken() {
            if (publicToken) {
                try {
                    const accessToken = await getAccessToken();
                    const response = await fetch('http://localhost:8080/exchange', {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                            'Authorization': `Bearer ${accessToken}`
                        },
                        body: JSON.stringify({public_token: publicToken})
                    });
                    const result = await response.json();
                    console.log('Token exchanged successfully:', result);
                    setExchangeResult(result);
                } catch (error) {
                    console.error('Error exchanging token:', error);
                }
            }
        }

        exchangeToken();
    }, [publicToken]);      
    return (
        <DashboardContent />
    )
}