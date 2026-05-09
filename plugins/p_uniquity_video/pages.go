package p_uniquity_video

import (
	"context"
	"fmt"
	"strings"

	"github.com/UniquityVentures/lago/components"
	"github.com/UniquityVentures/lago/getters"
	"github.com/UniquityVentures/lago/lago"
	"github.com/UniquityVentures/lago/plugins/p_filesystem"
	uniqempl "github.com/UniquityVentures/uniquity/plugins/p_uniquity_employees"
	"gorm.io/gorm"
)

// youtubeAPIFieldGetter returns context text for YouTube API fields, or "—" when unset.
func youtubeAPIFieldGetter(ctxKey string) getters.Getter[string] {
	return func(ctx context.Context) (string, error) {
		s, _ := getters.Key[string](ctxKey)(ctx)
		if strings.TrimSpace(s) == "" {
			return "—", nil
		}
		return s, nil
	}
}

// assignedToEmployeeGetter loads the employee (and User) for the Assigned to FK picker label.
func assignedToEmployeeGetter(ctx context.Context) (uniqempl.Employee, error) {
	var zero uniqempl.Employee
	fk, err := getters.Key[uint]("$in.AssignedToID")(ctx)
	if err != nil {
		return zero, err
	}
	if fk == 0 {
		return zero, nil
	}
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return zero, err
	}
	return gorm.G[uniqempl.Employee](db).Preload("User", nil).Where("id = ?", fk).Take(ctx)
}

// editedVideoForPublishedFKGetter loads EditedVideo with RawFootage for the published-video
// form’s InputForeignKey display ($in.RawFootage.Title). Plain [getters.Association] does not preload RawFootage.
func editedVideoForPublishedFKGetter(ctx context.Context) (EditedVideo, error) {
	var zero EditedVideo
	id, err := getters.Key[uint]("$in.EditedVideoID")(ctx)
	if err != nil {
		return zero, err
	}
	if id == 0 {
		return zero, nil
	}
	db, err := getters.DBFromContext(ctx)
	if err != nil {
		return zero, err
	}
	return gorm.G[EditedVideo](db).Preload("RawFootage", nil).Where("id = ?", id).Take(ctx)
}

// publishedVideoEditorPointsSeed pre-fills ToEmployeeID for the award form from the raw-footage assignee.
func publishedVideoEditorPointsSeed(ctx context.Context) (uniqempl.PointsTransaction, error) {
	var zero uniqempl.PointsTransaction
	pv, err := getters.Key[PublishedVideo]("publishedVideo")(ctx)
	if err != nil {
		return zero, err
	}
	id := pv.EditedVideo.RawFootage.AssignedToID
	if id == 0 {
		return zero, fmt.Errorf("no responsible editor: raw footage has no assignee")
	}
	return uniqempl.PointsTransaction{ToEmployeeID: id}, nil
}

func init() {
	registerVideoMenu()
	registerHubPage()
	registerRawPages()
	registerEditedPages()
	registerPublishedPages()
}

func registerVideoMenu() {
	lago.RegistryPage.Register("video.MainMenu", &components.SidebarMenu{
		Title: getters.Static("Video editors"),
		Back: &components.SidebarMenuItem{
			Title: getters.Static("Back to Home"),
			Url:   lago.RoutePath("dashboard.AppsPage", nil),
		},
		Children: []components.PageInterface{
			&components.SidebarMenuItem{Title: getters.Static("Overview"), Url: lago.RoutePath("video.DefaultRoute", nil), Icon: "home"},
			&components.SidebarMenuItem{Title: getters.Static("Raw footage"), Url: lago.RoutePath("video.RawListRoute", nil), Icon: "folder"},
			&components.SidebarMenuItem{Title: getters.Static("Edited videos"), Url: lago.RoutePath("video.EditedListRoute", nil), Icon: "scissors"},
			&components.SidebarMenuItem{Title: getters.Static("Published"), Url: lago.RoutePath("video.PublishedListRoute", nil), Icon: "play"},
		},
	})

	lago.RegistryPage.Register("video.RawDetailMenu", detailMenu("rawFootage", "Raw footage", "video.RawListRoute", "video.RawDetailRoute", "video.RawUpdateRoute"))
	lago.RegistryPage.Register("video.EditedDetailMenu", detailMenu("editedVideo", "Edited video", "video.EditedListRoute", "video.EditedDetailRoute", "video.EditedUpdateRoute"))
	lago.RegistryPage.Register("video.PublishedDetailMenu", publishedDetailMenu())
}

// publishedDetailMenu is like [detailMenu] with an extra superuser-only action to award points
// to the raw-footage assignee (responsible editor).
func publishedDetailMenu() *components.SidebarMenu {
	idGetter := getters.Key[uint]("publishedVideo.ID")
	return &components.SidebarMenu{
		Title: getters.Format("%s #%d", getters.Any(getters.Static("Published video")), getters.Any(idGetter)),
		Back: &components.SidebarMenuItem{
			Title: getters.Static("Back to list"),
			Url:   lago.RoutePath("video.PublishedListRoute", nil),
		},
		Children: []components.PageInterface{
			&components.SidebarMenuItem{
				Title: getters.Static("Detail"),
				Url: lago.RoutePath("video.PublishedDetailRoute", map[string]getters.Getter[any]{
					"id": getters.Any(idGetter),
				}),
			},
			&components.SidebarMenuItem{
				Title: getters.Static("Edit"),
				Url: lago.RoutePath("video.PublishedUpdateRoute", map[string]getters.Getter[any]{
					"id": getters.Any(idGetter),
				}),
			},
			&components.SidebarMenuItem{
				Page:  components.Page{Roles: []string{"superuser"}},
				Title: getters.Static("Give points to editor"),
				Icon:  "currency-dollar",
				Url: lago.RoutePath("video.PublishedEditorPointsCreateRoute", map[string]getters.Getter[any]{
					"id": getters.Any(idGetter),
				}),
			},
		},
	}
}

func detailMenu(ctxKey, title, listRoute, detailRoute, updateRoute string) *components.SidebarMenu {
	return &components.SidebarMenu{
		Title: getters.Format("%s #%d", getters.Any(getters.Static(title)), getters.Any(getters.Key[uint](ctxKey+".ID"))),
		Back: &components.SidebarMenuItem{
			Title: getters.Static("Back to list"),
			Url:   lago.RoutePath(listRoute, nil),
		},
		Children: []components.PageInterface{
			&components.SidebarMenuItem{
				Title: getters.Static("Detail"),
				Url:   lago.RoutePath(detailRoute, map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint](ctxKey + ".ID"))}),
			},
			&components.SidebarMenuItem{
				Title: getters.Static("Edit"),
				Url:   lago.RoutePath(updateRoute, map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint](ctxKey + ".ID"))}),
			},
		},
	}
}

func registerHubPage() {
	lago.RegistryPage.Register("video.HubPage", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.MainMenu"}},
		Children: []components.PageInterface{
			&components.ContainerColumn{
				Classes: "p-6 max-w-2xl",
				Children: []components.PageInterface{
					&components.FieldTitle{Getter: getters.Static("Video pipeline")},
					&components.FieldText{Getter: getters.Static("Manage raw footage, edited outputs, and YouTube publications from the sidebar.")},
				},
			},
		},
	})
}

func registerRawPages() {
	createN := getters.Static("video.RawFootageCreateForm")
	updateN := getters.Static("video.RawFootageUpdateForm")
	deleteN := getters.Static("video.RawFootageDeleteForm")

	formInputs := []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.Title"),
			Children: []components.PageInterface{
				&components.InputText{Label: "Title", Name: "Title", Required: true, Getter: getters.Key[string]("$in.Title")},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.Files"),
			Children: []components.PageInterface{
				&p_filesystem.InputMultiVNode{
					Label: "Files", Name: "Files", Required: false, Classes: "w-full",
					VNode: func(ctx context.Context) ([]p_filesystem.VNode, error) {
						if nodes, err := getters.Key[[]p_filesystem.VNode]("$in.Files")(ctx); err == nil && len(nodes) > 0 {
							return nodes, nil
						}
						return getters.AssociationList[p_filesystem.VNode](
							getters.AssociationIDs(getters.ContextKeyIn, "Files"),
							"",
						)(ctx)
					},
					AllowedFiletypes: []string{".mp4", ".mov", ".webm", ".mkv", ".mxf", ".avi", ".m4v", ".pdf", ".jpg", ".jpeg", ".png", ".webp"},
					Path: getters.Static("/video/raw-footage/uploads/"),
				},
			},
		},
		&components.InputForeignKey[uniqempl.Employee]{
			Name: "AssignedToID", Label: "Assigned to", Required: true,
			Url: lago.RoutePath("video.EmployeeSelectRoute", nil),
			Display: getters.Key[string]("$in.User.Name"), Placeholder: "Select employee…",
			Getter: assignedToEmployeeGetter,
		},
	}

	lago.RegistryPage.Register("video.RawFootageTable", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.MainMenu"}},
		Children: []components.PageInterface{
			&components.DataTable[RawFootage]{
				UID: "raw-footage-table", Classes: "w-full",
				Data: getters.Key[components.ObjectList[RawFootage]]("rawFootages"),
				Actions: []components.PageInterface{
					&components.TableButtonCreate{Link: lago.RoutePath("video.RawCreateRoute", nil)},
				},
				RowAttr: getters.RowAttrNavigate(lago.RoutePath("video.RawDetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))})),
				Columns: []components.TableColumn{
					{Label: "Title", Name: "Title", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$row.Title")}}},
					{Label: "Assigned to", Name: "AssignedTo", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$row.AssignedTo.User.Name")}}},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.RawFootageCreateForm", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.MainMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name: createN, ActionURL: lago.RoutePath("video.RawCreateRoute", nil),
				Children: []components.PageInterface{
					&components.FormComponent[RawFootage]{
						Attr: getters.FormBubbling(createN), Title: "New raw footage",
						ChildrenInput:  formInputs,
						ChildrenAction: []components.PageInterface{&components.ButtonSubmit{Label: "Save"}},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.RawFootageUpdateForm", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.RawDetailMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name: updateN,
				ActionURL: lago.RoutePath("video.RawUpdateRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("rawFootage.ID")),
				}),
				Children: []components.PageInterface{
					&components.FormComponent[RawFootage]{
						Getter: getters.Key[RawFootage]("rawFootage"), Attr: getters.FormBubbling(updateN),
						Title: "Edit raw footage", ChildrenInput: formInputs,
						ChildrenAction: []components.PageInterface{
							&components.ContainerRow{
								Classes: "flex flex-wrap justify-end gap-2",
								Children: []components.PageInterface{
									&components.ButtonSubmit{Label: "Update"},
									&components.ButtonModalForm{
										Label: "Delete", Icon: "trash", Name: deleteN,
										Url: lago.RoutePath("video.RawDeleteRoute", map[string]getters.Getter[any]{
											"id": getters.Any(getters.Key[uint]("rawFootage.ID")),
										}),
										FormPostURL: lago.RoutePath("video.RawDeleteRoute", map[string]getters.Getter[any]{
											"id": getters.Any(getters.Key[uint]("rawFootage.ID")),
										}),
										ModalUID: "raw-delete-modal", Classes: "btn-error",
									},
								},
							},
						},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.RawFootageDeleteForm", &components.Modal{
		UID: "raw-delete-modal",
		Children: []components.PageInterface{
			&components.DeleteConfirmation{
				Title: "Delete raw footage?", Message: "This cannot be undone.",
				Attr: getters.FormBubbling(getters.Key[string]("$get.name")),
			},
		},
	})

	lago.RegistryPage.Register("video.RawFootageDetail", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.RawDetailMenu"}},
		Children: []components.PageInterface{
			&components.Detail[RawFootage]{
				Getter: getters.Key[RawFootage]("rawFootage"),
				Children: []components.PageInterface{
					&components.ContainerColumn{Classes: "p-4 gap-2",
						Children: []components.PageInterface{
							&components.LabelInline{Title: "Title", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$in.Title")}}},
							&components.LabelInline{Title: "Files", Children: []components.PageInterface{
								&components.FieldList[p_filesystem.VNode]{
									Getter:  getters.Key[[]p_filesystem.VNode]("$in.Files"),
									Classes: "flex flex-col gap-2",
									Children: []components.PageInterface{
										&p_filesystem.FieldFile{VNode: getters.Key[p_filesystem.VNode]("$row")},
									},
								},
							}},
							&components.LabelInline{Title: "Assigned to", Children: []components.PageInterface{
								&components.InputForeignKey[uniqempl.Employee]{
									Name: "AssignedToID", Label: "", Required: true,
									Classes:     "pointer-events-none",
									Url:         lago.RoutePath("video.EmployeeSelectRoute", nil),
									Display:     getters.Key[string]("$in.User.Name"),
									Placeholder: "—",
									Getter:      assignedToEmployeeGetter,
								},
							}},
						}},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.RawFootageSelectionTable", &components.Modal{
		UID: "raw-footage-select-modal",
		Children: []components.PageInterface{
			&components.DataTable[RawFootage]{
				UID: "raw-footage-select-table", Title: "Select raw footage",
				Data:    getters.Key[components.ObjectList[RawFootage]]("rawFootages"),
				RowAttr: getters.RowAttrSelect("RawFootageID", getters.Key[uint]("$row.ID"), getters.Key[string]("$row.Title")),
				Columns: []components.TableColumn{
					{Label: "Title", Name: "Title", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$row.Title")}}},
				},
			},
		},
	})

	// Uses AssignedToID in RowAttrSelect so fk-select.name matches InputForeignKey.Name on raw footage forms.
	// employees.EmployeeSelectionTable hard-codes ToEmployeeID for the points form.
	lago.RegistryPage.Register("video.EmployeeSelectionTable", &components.Modal{
		UID: "video-employee-select-modal",
		Children: []components.PageInterface{
			&components.DataTable[uniqempl.Employee]{
				UID:   "video-employee-select-table",
				Title: "Select employee",
				Data:  getters.Key[components.ObjectList[uniqempl.Employee]]("employees"),
				RowAttr: getters.RowAttrSelect("AssignedToID", getters.Key[uint]("$row.ID"), getters.Key[string]("$row.User.Name")),
				Columns: []components.TableColumn{
					{Label: "User", Name: "User.Name", Children: []components.PageInterface{
						&components.FieldText{Getter: getters.Key[string]("$row.User.Name")},
					}},
					{Label: "Email", Name: "User.Email", Children: []components.PageInterface{
						&components.FieldText{Getter: getters.Key[string]("$row.User.Email")},
					}},
				},
			},
		},
	})
}

func registerEditedPages() {
	createN := getters.Static("video.EditedVideoCreateForm")
	updateN := getters.Static("video.EditedVideoUpdateForm")
	deleteN := getters.Static("video.EditedVideoDeleteForm")

	inputs := []components.PageInterface{
		&components.ContainerError{
			Error: getters.Key[error]("$error.RawFootageID"),
			Children: []components.PageInterface{
				&components.InputForeignKey[RawFootage]{
					Name: "RawFootageID", Label: "Raw footage", Required: true,
					Url: lago.RoutePath("video.RawSelectRoute", nil),
					Display: getters.Key[string]("$in.Title"), Placeholder: "Select raw footage…",
					Getter: getters.Association[RawFootage, uint](getters.Key[uint]("$in.RawFootageID")),
				},
			},
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.EditedVNodeID"),
			Children: []components.PageInterface{
				&p_filesystem.InputVNode{
					Label: "Output file", Name: "EditedVNodeID", Required: true, Classes: "w-full",
					VNode: func(ctx context.Context) (p_filesystem.VNode, error) {
						var zero p_filesystem.VNode
						if id, err := getters.Deref(getters.Key[*uint]("$in.EditedVNodeID"))(ctx); err == nil && id != 0 {
							return getters.Association[p_filesystem.VNode](getters.Static(id))(ctx)
						}
						if id, err := getters.Key[uint]("$in.EditedVNodeID")(ctx); err == nil && id != 0 {
							return getters.Association[p_filesystem.VNode](getters.Static(id))(ctx)
						}
						return zero, nil
					},
					AllowedFiletypes: []string{".mp4", ".mov", ".webm", ".mkv", ".mxf", ".avi", ".m4v", ".pdf"},
					Path: getters.Static("/video/edited/uploads/"),
				},
			},
		},
	}

	lago.RegistryPage.Register("video.EditedVideoTable", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.MainMenu"}},
		Children: []components.PageInterface{
			&components.DataTable[EditedVideo]{
				UID: "edited-video-table", Classes: "w-full",
				Data: getters.Key[components.ObjectList[EditedVideo]]("editedVideos"),
				Actions: []components.PageInterface{
					&components.TableButtonCreate{Link: lago.RoutePath("video.EditedCreateRoute", nil)},
				},
				RowAttr: getters.RowAttrNavigate(lago.RoutePath("video.EditedDetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))})),
				Columns: []components.TableColumn{
					{Label: "Raw title", Name: "Raw", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$row.RawFootage.Title")}}},
					{Label: "Output file", Name: "VNode", Children: []components.PageInterface{&p_filesystem.FieldFile{
						VNode: getters.Key[p_filesystem.VNode]("$row.EditedVNode"),
					}}},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.EditedVideoCreateForm", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.MainMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name: createN, ActionURL: lago.RoutePath("video.EditedCreateRoute", nil),
				Children: []components.PageInterface{
					&components.FormComponent[EditedVideo]{
						Attr: getters.FormBubbling(createN), Title: "New edited video",
						ChildrenInput:  inputs,
						ChildrenAction: []components.PageInterface{&components.ButtonSubmit{Label: "Save"}},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.EditedVideoUpdateForm", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.EditedDetailMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name: updateN,
				ActionURL: lago.RoutePath("video.EditedUpdateRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("editedVideo.ID")),
				}),
				Children: []components.PageInterface{
					&components.FormComponent[EditedVideo]{
						Getter: getters.Key[EditedVideo]("editedVideo"), Attr: getters.FormBubbling(updateN),
						Title: "Edit edited video", ChildrenInput: inputs,
						ChildrenAction: []components.PageInterface{
							&components.ContainerRow{
								Classes: "flex flex-wrap justify-end gap-2",
								Children: []components.PageInterface{
									&components.ButtonSubmit{Label: "Update"},
									&components.ButtonModalForm{
										Label: "Delete", Icon: "trash", Name: deleteN,
										Url: lago.RoutePath("video.EditedDeleteRoute", map[string]getters.Getter[any]{
											"id": getters.Any(getters.Key[uint]("editedVideo.ID")),
										}),
										FormPostURL: lago.RoutePath("video.EditedDeleteRoute", map[string]getters.Getter[any]{
											"id": getters.Any(getters.Key[uint]("editedVideo.ID")),
										}),
										ModalUID: "edited-delete-modal", Classes: "btn-error",
									},
								},
							},
						},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.EditedVideoDeleteForm", &components.Modal{
		UID: "edited-delete-modal",
		Children: []components.PageInterface{
			&components.DeleteConfirmation{
				Title: "Delete edited video?", Message: "This cannot be undone.",
				Attr:  getters.FormBubbling(getters.Key[string]("$get.name")),
			},
		},
	})

	lago.RegistryPage.Register("video.EditedVideoDetail", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.EditedDetailMenu"}},
		Children: []components.PageInterface{
			&components.Detail[EditedVideo]{
				Getter: getters.Key[EditedVideo]("editedVideo"),
				Children: []components.PageInterface{
					&components.ContainerColumn{Classes: "p-4 gap-2", Children: []components.PageInterface{
						&components.LabelInline{Title: "Raw footage", Children: []components.PageInterface{&components.FieldLink{
							Classes: "link link-hover",
							Href: lago.RoutePath("video.RawDetailRoute", map[string]getters.Getter[any]{
								"id": getters.Any(getters.Key[uint]("$in.RawFootageID")),
							}),
							Label: getters.Key[string]("$in.RawFootage.Title"),
						}}},
						&components.LabelInline{Title: "Assigned to", Children: []components.PageInterface{&components.FieldLink{
							Classes: "link link-hover",
							Href: lago.RoutePath("employees.EmployeeDetailRoute", map[string]getters.Getter[any]{
								"id": getters.Any(getters.Key[uint]("$in.RawFootage.AssignedToID")),
							}),
							Label: getters.Key[string]("$in.RawFootage.AssignedTo.User.Name"),
						}}},
						&components.LabelInline{Title: "Raw files", Children: []components.PageInterface{
							&components.FieldList[p_filesystem.VNode]{
								Getter:  getters.Key[[]p_filesystem.VNode]("$in.RawFootage.Files"),
								Classes: "flex flex-col gap-2",
								Children: []components.PageInterface{
									&p_filesystem.FieldFile{VNode: getters.Key[p_filesystem.VNode]("$row")},
								},
							},
						}},
						&components.LabelInline{Title: "Output file", Children: []components.PageInterface{&p_filesystem.FieldFile{
							VNode: getters.Key[p_filesystem.VNode]("$in.EditedVNode"),
						}}},
					}},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.EditedVideoSelectionTable", &components.Modal{
		UID: "edited-video-select-modal",
		Children: []components.PageInterface{
			&components.DataTable[EditedVideo]{
				UID: "edited-video-select-table", Title: "Select edited video",
				Data:    getters.Key[components.ObjectList[EditedVideo]]("editedVideos"),
				RowAttr: getters.RowAttrSelect("EditedVideoID", getters.Key[uint]("$row.ID"), getters.Key[string]("$row.RawFootage.Title")),
				Columns: []components.TableColumn{
					{Label: "Raw footage", Name: "Raw", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$row.RawFootage.Title")}}},
				},
			},
		},
	})
}

func registerPublishedPages() {
	createN := getters.Static("video.PublishedVideoCreateForm")
	updateN := getters.Static("video.PublishedVideoUpdateForm")
	deleteN := getters.Static("video.PublishedVideoDeleteForm")
	editorPointsN := getters.Static("video.PublishedEditorPointsForm")

	inputs := []components.PageInterface{
		&components.InputForeignKey[EditedVideo]{
			Name: "EditedVideoID", Label: "Edited video", Required: true,
			Url: lago.RoutePath("video.EditedSelectRoute", nil),
			Display: getters.Key[string]("$in.RawFootage.Title"), Placeholder: "Select edited cut…",
			Getter: editedVideoForPublishedFKGetter,
		},
		&components.ContainerError{
			Error: getters.Key[error]("$error.YouTubeVideoID"),
			Children: []components.PageInterface{
				&InputYouTubeVideoID{
					Label: "YouTube link or video ID", Name: "YouTubeVideoID", Required: true,
					Getter: getters.Key[string]("$in.YouTubeVideoID"),
				},
			},
		},
	}

	lago.RegistryPage.Register("video.PublishedVideoTable", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.MainMenu"}},
		Children: []components.PageInterface{
			&components.DataTable[PublishedVideo]{
				UID: "published-video-table", Classes: "w-full",
				Data: getters.Key[components.ObjectList[PublishedVideo]]("publishedVideos"),
				Actions: []components.PageInterface{
					&components.TableButtonCreate{Link: lago.RoutePath("video.PublishedCreateRoute", nil)},
				},
				RowAttr: getters.RowAttrNavigate(lago.RoutePath("video.PublishedDetailRoute", map[string]getters.Getter[any]{"id": getters.Any(getters.Key[uint]("$row.ID"))})),
				Columns: []components.TableColumn{
					{Label: "YouTube ID", Name: "YT", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$row.YouTubeVideoID")}}},
					{Label: "Raw title", Name: "Raw", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$row.EditedVideo.RawFootage.Title")}}},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.PublishedVideoCreateForm", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.MainMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name: createN, ActionURL: lago.RoutePath("video.PublishedCreateRoute", nil),
				Children: []components.PageInterface{
					&components.FormComponent[PublishedVideo]{
						Attr: getters.FormBubbling(createN), Title: "New published video",
						ChildrenInput:  inputs,
						ChildrenAction: []components.PageInterface{&components.ButtonSubmit{Label: "Save"}},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.PublishedVideoUpdateForm", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.PublishedDetailMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name: updateN,
				ActionURL: lago.RoutePath("video.PublishedUpdateRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("publishedVideo.ID")),
				}),
				Children: []components.PageInterface{
					&components.FormComponent[PublishedVideo]{
						Getter: getters.Key[PublishedVideo]("publishedVideo"), Attr: getters.FormBubbling(updateN),
						Title: "Edit published video", ChildrenInput: inputs,
						ChildrenAction: []components.PageInterface{
							&components.ContainerRow{
								Classes: "flex flex-wrap justify-end gap-2",
								Children: []components.PageInterface{
									&components.ButtonSubmit{Label: "Update"},
									&components.ButtonModalForm{
										Label: "Delete", Icon: "trash", Name: deleteN,
										Url: lago.RoutePath("video.PublishedDeleteRoute", map[string]getters.Getter[any]{
											"id": getters.Any(getters.Key[uint]("publishedVideo.ID")),
										}),
										FormPostURL: lago.RoutePath("video.PublishedDeleteRoute", map[string]getters.Getter[any]{
											"id": getters.Any(getters.Key[uint]("publishedVideo.ID")),
										}),
										ModalUID: "published-delete-modal", Classes: "btn-error",
									},
								},
							},
						},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.PublishedVideoDeleteForm", &components.Modal{
		UID: "published-delete-modal",
		Children: []components.PageInterface{
			&components.DeleteConfirmation{
				Title: "Delete published video?", Message: "This cannot be undone.",
				Attr:  getters.FormBubbling(getters.Key[string]("$get.name")),
			},
		},
	})

	lago.RegistryPage.Register("video.PublishedVideoDetail", &components.ShellScaffold{
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.PublishedDetailMenu"}},
		Children: []components.PageInterface{
			&components.Detail[PublishedVideo]{
				Getter: getters.Key[PublishedVideo]("publishedVideo"),
				Children: []components.PageInterface{
					&components.ContainerColumn{Classes: "p-4 gap-2", Children: []components.PageInterface{
						&components.LabelInline{Title: "YouTube video", Children: []components.PageInterface{&components.FieldLink{
							Classes: "link link-hover break-all",
							Href: func(ctx context.Context) (string, error) {
								id, err := getters.Key[string]("$in.YouTubeVideoID")(ctx)
								if err != nil {
									return "", err
								}
								u := YouTubeWatchURL(id)
								if u == "" {
									return "", nil
								}
								return u, nil
							},
							Label: getters.Key[string]("$in.YouTubeVideoID"),
						}}},
						&components.LabelInline{Title: "YouTube Studio", Children: []components.PageInterface{&components.FieldLink{
							Classes: "link link-hover",
							Href: func(ctx context.Context) (string, error) {
								id, err := getters.Key[string]("$in.YouTubeVideoID")(ctx)
								if err != nil {
									return "", err
								}
								u := YouTubeStudioVideoURL(id)
								if u == "" {
									return "", nil
								}
								return u, nil
							},
							Label: getters.Static("Open video in YouTube Studio"),
						}}},
						&components.LabelInline{Title: "Title (YouTube)", Children: []components.PageInterface{
							&components.FieldText{Getter: youtubeAPIFieldGetter(ctxYouTubeSnippetTitle)},
						}},
						&components.LabelInline{Title: "Published on YouTube", Children: []components.PageInterface{
							&components.FieldText{Getter: youtubeAPIFieldGetter(ctxYouTubePublishedAtDisplay)},
						}},
						&components.LabelInline{Title: "YouTube upload status", Children: []components.PageInterface{
							&components.FieldText{Getter: youtubeAPIFieldGetter(ctxYouTubeUploadStatus)},
						}},
						&components.LabelInline{Title: "Views", Children: []components.PageInterface{
							&components.FieldText{Getter: youtubeAPIFieldGetter(ctxYouTubeViewCount)},
						}},
						&components.LabelInline{Title: "Likes", Children: []components.PageInterface{
							&components.FieldText{Getter: youtubeAPIFieldGetter(ctxYouTubeLikeCount)},
						}},
						&components.LabelInline{Title: "Comments", Children: []components.PageInterface{
							&components.FieldText{Getter: youtubeAPIFieldGetter(ctxYouTubeCommentCount)},
						}},
						&components.LabelInline{Title: "Edited from (raw)", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$in.EditedVideo.RawFootage.Title")}}},
						&components.LabelInline{Title: "Assigned to", Children: []components.PageInterface{&components.FieldLink{
							Classes: "link link-hover",
							Href: lago.RoutePath("employees.EmployeeDetailRoute", map[string]getters.Getter[any]{
								"id": getters.Any(getters.Key[uint]("$in.EditedVideo.RawFootage.AssignedToID")),
							}),
							Label: getters.Key[string]("$in.EditedVideo.RawFootage.AssignedTo.User.Name"),
						}}},
					}},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.PublishedEditorPointsForm", &components.ShellScaffold{
		Page:    components.Page{Roles: []string{"superuser"}},
		Sidebar: []components.PageInterface{lago.DynamicPage{Name: "video.PublishedDetailMenu"}},
		Children: []components.PageInterface{
			&components.FormListenBoostedPost{
				Name: editorPointsN,
				ActionURL: lago.RoutePath("video.PublishedEditorPointsCreateRoute", map[string]getters.Getter[any]{
					"id": getters.Any(getters.Key[uint]("publishedVideo.ID")),
				}),
				Children: []components.PageInterface{
					&components.FormComponent[uniqempl.PointsTransaction]{
						Getter:   publishedVideoEditorPointsSeed,
						Attr:     getters.FormBubbling(editorPointsN),
						Title:    "Give points to editor",
						Subtitle: "Awards points to the employee assigned to the source raw footage. The granting user is recorded as the signed-in superuser.",
						ChildrenInput: []components.PageInterface{
							&components.LabelInline{
								Title: "Responsible editor",
								Children: []components.PageInterface{
									&components.FieldText{Getter: getters.Key[string]("publishedVideo.EditedVideo.RawFootage.AssignedTo.User.Name")},
								},
							},
							&components.InputForeignKey[uniqempl.Employee]{
								Name:     "ToEmployeeID",
								Hidden:   true,
								Required: true,
								Getter:   getters.Association[uniqempl.Employee, uint](getters.Key[uint]("$in.ToEmployeeID")),
							},
							&uniqempl.InputPointsDecimal{
								Label:    "Points",
								Name:     "Points",
								Required: true,
								Getter:   getters.Key[uniqempl.PointsDecimal]("$in.Points"),
							},
						},
						ChildrenAction: []components.PageInterface{
							&components.ButtonSubmit{Label: "Award points"},
						},
					},
				},
			},
		},
	})

	lago.RegistryPage.Register("video.PublishedVideoSelectionTable", &components.Modal{
		UID: "published-video-select-modal",
		Children: []components.PageInterface{
			&components.DataTable[PublishedVideo]{
				UID: "published-video-select-table", Title: "Select published video",
				Data:    getters.Key[components.ObjectList[PublishedVideo]]("publishedVideos"),
				RowAttr: getters.RowAttrSelect("PublishedVideoID", getters.Key[uint]("$row.ID"), getters.Key[string]("$row.YouTubeVideoID")),
				Columns: []components.TableColumn{
					{Label: "YouTube ID", Name: "YT", Children: []components.PageInterface{&components.FieldText{Getter: getters.Key[string]("$row.YouTubeVideoID")}}},
				},
			},
		},
	})
}
