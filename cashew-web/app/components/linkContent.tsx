'use client'
import { useRouter } from 'next/router';
import { useState, useEffect, useCallback } from 'react';
import { PlaidLinkOnSuccess, PlaidLinkOnSuccessMetadata, usePlaidLink } from "react-plaid-link";

export function LinkContent({initialLinkToken, accessToken}: { initialLinkToken: string, accessToken: string}) {
    const router = useRouter()
    const [linkToken, setLinkToken] = useState(initialLinkToken);
    const onSuccess = useCallback<PlaidLinkOnSuccess>(
        (public_token: string, metadata: PlaidLinkOnSuccessMetadata) => {
            console.log('success:', public_token, metadata)
            router.push(`/dashboard?public_token=${encodeURIComponent(public_token)}`)
        },
        [],      
    )

    const { open, ready } = usePlaidLink({
        token: linkToken,
        onSuccess: onSuccess
    });

    return (
        <div>
            <h1>Cashew Link</h1>
            {linkToken && <pre>{linkToken}</pre>}
            <button onClick={() => open()} disabled={!ready}>
                Connect a bank account
            </button>
            <a href="/api/auth/logout">Logout</a>
        </div>
    );
}