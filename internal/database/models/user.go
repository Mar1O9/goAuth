// Copyright 2025 Mahmoud Abdelrahman <mahmud.yousif.04@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


package models

import (
	"time"

	"gorm.io/gorm"
)

// System User model
type User struct {
	ID           uint           `gorm:"primaryKey"`
	Name         string         `gorm:"size:250;"`
	Username     string         `gorm:"uniqueIndex;size:100;not null"`
	Email        string         `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash []byte         `gorm:"not null"`
	IsStaff      bool           `gorm:"default:false"` // For staff members
	IsSuperuser  bool           `gorm:"default:false"` // For admins
	IsActive     bool           `gorm:"default:true"`  // Can be banned or active
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoCreateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
