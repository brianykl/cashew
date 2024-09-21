

export default function Home() {
  return (
    <>
      <div className="flex items-center justify-center  h-[60vh]">
        <h1 className="text-center text-4xl font-bold -mt-16">cashew</h1>
      </div>
      <div>
        <a href="/api/auth/login">Login</a>
        <a href="/api/auth/logout">Logout</a>
      </div>
    </>
  );
}
