'use client'
import { useRouter } from 'next/navigation';
import { useState, useEffect } from 'react';
import { PlaidLinkOnSuccess, PlaidLinkOnSuccessMetadata, usePlaidLink } from "react-plaid-link";


interface DashboardContentProps {
    initialLinkToken: string;
    accessToken: string;
    userId: string;
}

interface Account {
    id: string;
    name: string;
    type: string;
    subtype: string;
}

export function DashboardContent({initialLinkToken, accessToken, userId}: DashboardContentProps) {
    const [linkToken] = useState(initialLinkToken);
    const [accounts, setAccounts] = useState<Account[]>([]);

    const onSuccess: PlaidLinkOnSuccess = async (public_token, metadata) => {
        console.log('Plaid Link success:', public_token, metadata);
        try {
            const response = await fetch('http://localhost:8080/protected/exchange', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify({
                    public_token: public_token,
                    user_id: userId
                })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const result = await response.json();
            console.log("exchange result:", result);
        } catch (error) {
            console.error("Error exchanging token:", error);
        }
    };

    const config = {
        token: linkToken,
        onSuccess: onSuccess,
    };

    const { open, ready, error } = usePlaidLink(config);


    useEffect(() => {
        // Fetch connected accounts when component mounts
        fetchAccounts();
    }, []);

    const fetchAccounts = async () => {
        // Implement API call to fetch accounts
        // This is a placeholder. Replace with your actual API call.
        const mockAccounts: Account[] = [
            { id: '1', name: 'Checking Account', type: 'depository', subtype: 'checking' },
            { id: '2', name: 'Savings Account', type: 'depository', subtype: 'savings' },
            { id: '3', name: 'Credit Card', type: 'credit', subtype: 'credit card' },
        ];
        setAccounts(mockAccounts);
    };

    return (
        <div className="flex h-screen bg-gray-100">
            {/* Left sidebar with accounts table */}
            <div className="w-1/4 bg-white shadow-md p-6 overflow-y-auto">
                <h2 className="text-xl font-bold mb-4">Connected Accounts</h2>
                <table className="w-full">
                    <thead>
                        <tr>
                            <th className="text-left pb-2">Name</th>
                            <th className="text-left pb-2">Type</th>
                        </tr>
                    </thead>
                    <tbody>
                        {accounts.map((account) => (
                            <tr key={account.id} className="border-t">
                                <td className="py-2">{account.name}</td>
                                <td className="py-2">{account.subtype}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {/* Main content area */}
            <div className="flex-1 p-10">
                <h1 className="text-3xl font-bold mb-6">Cashew Dashboard</h1>
                <button 
                    onClick={() => open()} 
                    disabled={!ready}
                    className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                >
                    Connect a bank account
                </button>
                <a href="/api/auth/logout" className="ml-4 text-blue-500 hover:text-blue-700">Logout</a>
            </div>
        </div>
    );
}