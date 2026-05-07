package p_uniquity_video

import (
	"github.com/UniquityVentures/lago/lago"
	"github.com/UniquityVentures/lago/plugins/p_filesystem"
	uniqempl "github.com/UniquityVentures/uniquity/plugins/p_uniquity_employees"
	"gorm.io/gorm"
)

// RawFootage is source material for the video pipeline.
type RawFootage struct {
	gorm.Model

	Title string `gorm:"not null"`
	Files []p_filesystem.VNode `gorm:"many2many:raw_footage_files;"`

	AssignedToID uint            `gorm:"not null"`
	AssignedTo   uniqempl.Employee `gorm:"foreignKey:AssignedToID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// EditedVideo is a rendered cut linked to raw footage and an output file node.
type EditedVideo struct {
	gorm.Model

	RawFootageID uint       `gorm:"not null"`
	RawFootage   RawFootage `gorm:"foreignKey:RawFootageID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	EditedVNodeID uint               `gorm:"not null"`
	EditedVNode   p_filesystem.VNode `gorm:"foreignKey:EditedVNodeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// PublishedVideo is a YouTube publication of an edited video.
type PublishedVideo struct {
	gorm.Model

	EditedVideoID uint        `gorm:"not null"`
	EditedVideo   EditedVideo `gorm:"foreignKey:EditedVideoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	YouTubeVideoID string `gorm:"not null;size:32"`
}

func init() {
	lago.OnDBInit("p_uniquity_video.models", func(d *gorm.DB) *gorm.DB {
		lago.RegisterModel[RawFootage](d)
		lago.RegisterModel[EditedVideo](d)
		lago.RegisterModel[PublishedVideo](d)
		return d
	})

	lago.RegistryAdmin.Register("p_uniquity_video_raw", lago.AdminPanel[RawFootage]{
		SearchField: "Title",
		ListFields:  []string{"Title", "AssignedTo.User.Name", "UpdatedAt"},
		Preload:     []string{"Files", "AssignedTo", "AssignedTo.User"},
	})

	lago.RegistryAdmin.Register("p_uniquity_video_edited", lago.AdminPanel[EditedVideo]{
		SearchField: "RawFootage.Title",
		ListFields:  []string{"RawFootage.Title", "EditedVNode.Name", "UpdatedAt"},
		Preload:     []string{"RawFootage", "EditedVNode"},
	})

	lago.RegistryAdmin.Register("p_uniquity_video_published", lago.AdminPanel[PublishedVideo]{
		SearchField: "YouTubeVideoID",
		ListFields:  []string{"YouTubeVideoID", "EditedVideo.RawFootage.Title", "UpdatedAt"},
		Preload:     []string{"EditedVideo", "EditedVideo.RawFootage"},
	})
}
