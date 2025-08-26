"use client";
import { logout } from "@/auth/server/action";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";

export default function Logout() {
    const router = useRouter();

  const handleLogout = async () => {
    await logout();
    router.refresh();
  };

  return <Button onClick={handleLogout}>Logout</Button>;
}
