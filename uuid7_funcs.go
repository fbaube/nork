package nork

import (
    "uuid"
)

func (p *Nork)  UUID() uuid.UUID  { return p.UUIDv7 }
func (p *Nork) TableName() string { return string(p.UUIDv7[13:16]) }

