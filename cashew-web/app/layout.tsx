import { UserProvider } from "@auth0/nextjs-auth0/client";
import "./globals.css";
import Image from "next/image";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="relative min-h-screen">
        <div className="flex flex-col">
          <header className="flex justify-between items-center p-4">
            <div className="flex items-center space-x-2">
              <Image
                src="/logos/cashew.png"
                alt="Logo"
                width={100}
                height={100}
              />
              <span className="text-6xl font-bold">cashew</span>
            </div>
            <a
              href="/api/auth/logout"
              className="bg-blue-500 hover:bg-blue-700 text-white text-md font-bold py-0.5 px-6 mr-2 rounded"
            >
              logout
            </a>
          </header>

          <div className="flex-1">
            <UserProvider>
              <main className="p-6">{children}</main>
            </UserProvider>
          </div>
        </div>
      </body>
    </html>
  );
}
