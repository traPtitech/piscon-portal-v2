import { traQMarkdownIt, type Store } from '@traptitech/traq-markdown-it'

const store: Store = {
  getUser: () => undefined,
  getChannel: () => undefined,
  getUserGroup: () => undefined,
  getMe: () => undefined,
  getStampByName: () => undefined,
  getUserByName: () => undefined,
  generateUserHref: () => '',
  generateUserGroupHref: () => '',
  generateChannelHref: () => '',
}
const it = new traQMarkdownIt(store, undefined, 'https://q.trap.jp')

export const renderMarkdown = (markdown: string) => {
  return it.render(markdown)
}
