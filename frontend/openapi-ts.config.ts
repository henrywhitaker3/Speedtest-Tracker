import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: 'http://localhost:8080/docs/schema.yaml',
  output: 'app/client'
})
