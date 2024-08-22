'use client'
import { useState, useEffect } from 'react';
import { usePlaidLink } from "react-plaid-link";

export function DashboardContent({initialLinkToken}: { initialLinkToken: string }) {
    const [linkToken, setLinkToken] = useState(initialLinkToken);
    const { open, ready } = usePlaidLink({
        token: linkToken,
        onSuccess: (public_token, metadata) => {
            console.log('Success:', public_token, metadata);
            // Handle success
        },
    });

    return (
        <div>
            <h1>Cashew Dashboard</h1>
            {linkToken && <pre>{linkToken}</pre>}
            <button onClick={() => open()} disabled={!ready}>
                Connect a bank account
            </button>
            <a href="/api/auth/logout">Logout</a>
        </div>
    );
}