interface Account {
  name: string;
  type: string;
  available_balance: number;
}

export async function refreshAccounts(
  userId: string,
  accessToken: string
): Promise<Account[]> {
  try {
    const response = await fetch("http://localhost:8080/protected/accounts", {
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

    const accounts: Account[] = await response.json();
    return accounts;
  } catch (error) {
    console.error("Error fetching accounts:", error);
    throw error
  }
}
