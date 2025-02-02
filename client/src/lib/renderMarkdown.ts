import { micromark } from 'micromark'

export const renderMarkdown = (markdown: string) => {
  return micromark(markdown)
}
