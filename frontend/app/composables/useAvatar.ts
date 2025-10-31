export const useAvatar = () => {
  const getAvatarUrl = (
    avatarUrl: string | null | undefined,
    domain: string | null | undefined
  ) => {
    const fallbackDomain = domain || 'default'
    return (
      avatarUrl ||
      `https://api.dicebear.com/9.x/glass/svg?seed=${fallbackDomain}`
    )
  }
  return { getAvatarUrl }
}
