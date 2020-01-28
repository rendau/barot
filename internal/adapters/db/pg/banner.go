package pg

import (
	"context"
	"github.com/rendau/barot/internal/domain/entities"
)

func (d *St) BannerCreate(ctx context.Context, obj *entities.Banner) error {
	_, err := d.db.ExecContext(ctx, `
		insert into banner(id, slot_id, note)
		values($1, $2, $3)
	`, obj.ID, obj.SlotID, obj.Note)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return err
	}

	return nil
}

func (d *St) BannerList(ctx context.Context, pars entities.BannerFilterPars) ([]*entities.Banner, error) {
	args := map[string]interface{}{}

	if pars.ID != nil {
		args["id"] = *pars.ID
	}
	if pars.SlotID != nil {
		args["slot_id"] = *pars.SlotID
	}

	qWhere := ``
	for k := range args {
		if qWhere != `` {
			qWhere += ` and`
		}
		qWhere += ` ` + k + ` = :` + k
	}

	stmt, err := d.db.PrepareNamedContext(ctx, `
		select b.id, b.slot_id, b.note
		from banner b
		where `+qWhere+`
		order by b.id
	`)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, args)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return nil, err
	}
	defer rows.Close()

	items := make([]*entities.Banner, 0)

	for rows.Next() {
		item := &entities.Banner{}
		err = rows.Scan(
			&item.ID,
			&item.SlotID,
			&item.Note,
		)
		if err != nil {
			d.lg.Errorw(ErrMsg, err)
			return nil, err
		}
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		d.lg.Errorw(ErrMsg, err)
		return nil, err
	}

	return items, nil
}

func (d *St) BannerDelete(ctx context.Context, pars entities.BannerFilterPars) error {
	var err error

	args := map[string]interface{}{}

	if pars.ID != nil {
		args["id"] = *pars.ID
	}
	if pars.SlotID != nil {
		args["slot_id"] = *pars.SlotID
	}

	if len(args) == 0 {
		return nil
	}

	qWhere := ``
	for k := range args {
		if qWhere != `` {
			qWhere += ` and`
		}
		qWhere += ` ` + k + ` = :` + k
	}

	stmt, err := d.db.PrepareNamedContext(ctx, `
		delete from banner
		where `+qWhere+`
	`)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, args)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return err
	}

	return nil
}
