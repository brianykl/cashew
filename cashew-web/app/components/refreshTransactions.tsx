
interface Transaction {
    date: string,
    account_name: string,
    merchant_name: string,
    currency: string,
    amount: number,
    primary_category: string,
    detailed_category: string
}

export async function refreshTransactions(
    userId: string,
    accessToken: string
  ): Promise<Transaction[]> {
    try {
      const response = await fetch("http://localhost:8080/protected/transactions", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          user_id: userId,
        }),
      });
  
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      const accounts: Transaction[] = await response.json();
      return accounts;
    } catch (error) {
      console.error("Error fetching transactions:", error);
      throw error
    }
  }