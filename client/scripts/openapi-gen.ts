import * as Bun from 'bun'
import * as yaml from 'yaml'

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const stripDiscriminator = (schema: any) => {
  Object.entries(schema).forEach(([key, value]) => {
    if (!(typeof value === 'object')) return
    if (key === 'discriminator') {
      delete schema[key]
    } else {
      stripDiscriminator(value)
    }
  })
}

const file = await Bun.file('../openapi/openapi.yml').text()
const parsed = yaml.parse(file)
stripDiscriminator(parsed)
const temp = Bun.env['FILE'] ?? '/tmp/openapi.yml'
await Bun.write(temp, yaml.stringify(parsed))
