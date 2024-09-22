export interface Transaction {
  date: string;
  account_name: string;
  merchant_name: string;
  currency: string;
  amount: number;
  primary_category: string;
  detailed_category: string;
}

interface APITransaction {
  UserId: string;
  AccountId: string;
  AccountName: string;
  Amount: string;
  Currency: string;
  AuthorizedDate: string;
  MerchantName: string;
  PaymentChannel: string;
  PrimaryCategory: string;
  DetailedCategory: string;
  ConfidenceLevel: string;
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

      const data = await response.json();
      const apiTransactions: APITransaction[] = data.transactions;

      const transactions: Transaction[] = apiTransactions.map(t => ({
          date: t.AuthorizedDate,
          account_name: t.AccountName,
          merchant_name: t.MerchantName,
          currency: t.Currency,
          amount: parseFloat(t.Amount),
          primary_category: t.PrimaryCategory,
          detailed_category: t.DetailedCategory
      }));

      return transactions;
  } catch (error) {
      console.error("Error fetching transactions:", error);
      throw error;
  }
} 