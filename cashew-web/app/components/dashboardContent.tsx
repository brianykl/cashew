"use client";
import { useRouter } from "next/navigation"
import { useState, useEffect, useRef } from "react"
import {
  PlaidLinkOnSuccess,
  PlaidLinkOnSuccessMetadata,
  usePlaidLink,
} from "react-plaid-link"
import { Account, refreshAccounts } from "./refreshAccounts"
import { Transaction, refreshTransactions } from "./refreshTransactions"

interface DashboardContentProps {
  initialLinkToken: string;
  accessToken: string;
  userId: string;
}

// interface Account {
//   name: string;
//   type: string;
//   available_balance: number;
// }


// interface Transaction {
//   date: string,
//   account_name: string,
//   merchant_name: string,
//   currency: string,
//   amount: number,
//   primary_category: string,
//   detailed_category: string
// }

export function DashboardContent({
  initialLinkToken,
  accessToken,
  userId,
}: DashboardContentProps) {
  const [linkToken] = useState(initialLinkToken)
  const [accounts, setAccounts] = useState<Account[]>([])
  const [transactions, setTransactions] = useState<Transaction[]>([])
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
      fetchAccounts()
      fetchTransactions()
      initialFetchDone.current = true;
    }
  }, [])
  const fetchAccounts = async () => {
    try {
      const fetchedAccounts = await refreshAccounts(userId, accessToken)
      setAccounts(fetchedAccounts)
    } catch (error) {
      console.error("error fetching accounts:", error);
    }
  }
  
  const fetchTransactions = async () => {{
    try {
      const fetchedTransactions = await refreshTransactions(userId, accessToken)
      setTransactions(fetchedTransactions)
    }catch (error) {
      console.error("error fetching transactions:", error);
    }
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
      <table className="w-full">
          <thead className="">
            <tr>
              <th className="text-left pb-2">date</th>
              <th className="text-left pb-2">account name</th>
              <th className="text-left pb-2">merchant</th>
              <th className="text-left pb-2">currency</th>
              <th className="text-left pb-2">amount</th>
              <th className="text-left pb-2">primary category</th>
              <th className="text-left pb-2">detailed category</th>
            </tr>
          </thead>
          <tbody>
            {transactions.map((transaction, index) => (
              <tr key={index} className="border-t">
                <td className="py-2">{transaction.date}</td>
                <td className="py-2">{transaction.account_name}</td>
                <td className="py-2">{transaction.merchant_name}</td>
                <td className="py-2">{transaction.currency}</td>
                <td className="py-2">{transaction.amount.toFixed(2)}</td>
                <td className="py-2">${transaction.primary_category}</td>
                <td className="py-2">${transaction.detailed_category}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
