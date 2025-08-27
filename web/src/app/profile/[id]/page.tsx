import { getAuthUser } from "@/auth/server/action"
import Profile from "./profile";

export default async function ProfilePage() {
    const auth = await getAuthUser();
    return <Profile auth={auth}/>
}