"use client";
import { useRouter } from "next/navigation"
import { useState, useEffect, useRef } from "react"
import {
  PlaidLinkOnSuccess,
  PlaidLinkOnSuccessMetadata,
  usePlaidLink,
} from "react-plaid-link"
import { refreshAccounts } from "./refreshAccounts"
import { refreshTransactions } from "./refreshTransactions"

interface DashboardContentProps {
  initialLinkToken: string;
  accessToken: string;
  userId: string;
}

interface Account {
  name: string;
  type: string;
  available_balance: number;
}

export function DashboardContent({
  initialLinkToken,
  accessToken,
  userId,
}: DashboardContentProps) {
  const [linkToken] = useState(initialLinkToken)
  const [accounts, setAccounts] = useState<Account[]>([])
  const initialFetchDone = useRef(false);

  const onSuccess: PlaidLinkOnSuccess = async (public_token, metadata) => {
    console.log("Plaid Link success:", public_token, metadata);
    try {
      const response = await fetch("http://localhost:8080/protected/exchange", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          public_token: public_token,
          user_id: userId,
        }),
      })

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
    if (!initialFetchDone.current) {
      fetchAccounts();
      initialFetchDone.current = true;
    }
  }, []);
  const fetchAccounts = async () => {
    try {
      const fetchedAccounts = await refreshAccounts(userId, accessToken);
      setAccounts(fetchedAccounts);
    } catch (error) {
      console.error("Error fetching accounts:", error);
    }
  }

  
  return (
    <div className="flex h-full bg-gray-100 border-1 border-black">
      {/* Left sidebar with accounts table */}
      <div className="w-1/4 shadow-md flex flex-col p-4 overflow-y-auto border-2 border-leaf bg-yellow-50">
        <h2 className="text-xl text-center font-bold">connected accounts</h2>
        <div className="flex space-x-2">
          <button
            onClick={() => open()}
            disabled={!ready}
            className="bg-blue-500 hover:bg-blue-700 text-white text-sm font-bold py-0.5 px-4 rounded"
          >
            connect a bank account
          </button>
          <button
            onClick={() => refreshAccounts(userId, accessToken)}
            disabled={!ready}
            className="bg-blue-500 hover:bg-blue-700 text-white text-sm font-bold py-0. 5 px-4 rounded"
          >
            refresh accounts
          </button>
        </div>
        <table className="w-full">
          <thead className="">
            <tr>
              <th className="text-left pb-2">name</th>
              <th className="text-left pb-2">type</th>
              <th className="text-left pb-2">balance</th>
            </tr>
          </thead>
          <tbody>
            {accounts.map((account, index) => (
              <tr key={index} className="border-t">
                <td className="py-2">{account.name}</td>
                <td className="py-2">{account.type}</td>
                <td className="py-2">${account.available_balance.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Main content area */}
      <div className="flex-1 p-10">
        <h1 className="text-3xl font-bold mb-6">cashew dashboard</h1>
      </div>
    </div>
  );
}
