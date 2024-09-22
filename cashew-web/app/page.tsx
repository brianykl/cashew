

export default function Home() {
  return (
    <>
      <div className="flex items-center justify-center h-[60vh]">
        <h1 className="text-center text-xl font-bold -mt-16">where's ya cash going</h1>
      </div>
      <div className="flex items-center justify-center">
        <a href="/api/auth/login">Login</a>
      </div>
    </>
  );
}
