import { getSession } from "@auth0/nextjs-auth0";
import { redirect } from "next/navigation";

export default async function Dashboard() {
    const session = await getSession();

    if(!session) {
        redirect('api/auth/login')
    }

    return (
      <div>cashew dashboard
        <a href="/api/auth/logout">Logout</a>
      </div>
    );
  }