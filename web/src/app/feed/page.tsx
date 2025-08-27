import { getAuthUser } from "@/auth/server/action"
import { PersonalizedFeed } from "@/components/feed/personalized-feed"

export default async function FeedPage() {
  const auth = await getAuthUser();

  if (!auth) {
    return (
      <div className="container max-w-2xl mx-auto py-8 px-4">
        <p>Please log in to view your feed.</p>
      </div>
    );
  }

  return (
    <div className="container max-w-2xl mx-auto py-8 px-4">
      <PersonalizedFeed auth={auth} />
    </div>
  )
}
