import { useEffect, useState } from 'react'

export const DEFAULT_POST_IMAGE =
  'https://images.unsplash.com/photo-1526498460520-4c246339dccb?auto=format&fit=crop&w=1200&q=80'

export default function PostImage({ src, alt }) {
  const [failed, setFailed] = useState(false)
  useEffect(() => {
    setFailed(false)
  }, [src])

  const imageSrc = !failed && src ? src : DEFAULT_POST_IMAGE

  return <img className="post-image" src={imageSrc} alt={alt} onError={() => setFailed(true)} loading="lazy" />
}
