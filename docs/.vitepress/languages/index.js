import cds from './cds.textMate.json' with {type: 'json'}
import dcl from './dcl.textMate.json' with {type: 'json'}

export default [
    { ...cds, aliases: ['cds'] },
    { ...dcl, aliases: ['dcl'] }
]