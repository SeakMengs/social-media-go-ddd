import { getAuthUser } from "@/auth/server/action";
import { Button } from "@/components/ui/button";
import moment from "moment";
import Link from "next/link";
import Logout from "./logout";

export default async function Home() {
  const auth = await getAuthUser();

  return (
    <div>
      <h1>Welcome to the Home Page</h1>
      {auth ? (
        <>
          <p>
            Hello, {auth.user.username}!<br />
            Session Expire at:{" "}
            {moment(auth.session.expireAt).format("MMMM Do YYYY, h:mm:ss a")}
          </p>
          <Logout />
        </>
      ) : (
        <p>
          Please log in to access more features.
          <Link href="/login" className="text-blue-500 underline">
            Go to Login
          </Link>
        </p>
      )}
    </div>
  );
}
