package pg

import (
	"context"
	"github.com/rendau/barot/internal/domain/entities"
)

func (d *St) BannerCreate(ctx context.Context, pars entities.BannerCreatePars) error {
	_, err := d.Db.ExecContext(ctx, `
		insert into banner(id, slot_id, note)
		values($1, $2, $3)
	`, pars.ID, pars.SlotID, pars.Note)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return err
	}

	return nil
}

func (d *St) BannerDelete(ctx context.Context, pars entities.BannerDeletePars) error {
	var err error

	_, err = d.Db.ExecContext(ctx, `
		delete from banner
		where id = $1 and slot_id = $2
	`, pars.ID, pars.SlotID)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return err
	}

	return nil
}

func (d *St) BannerList(ctx context.Context, pars entities.BannerListPars) ([]*entities.Banner, error) {
	rows, err := d.Db.QueryxContext(ctx, `
		select b.id,
		       b.slot_id,
		       coalesce(s.show_cnt, 0),
		       coalesce(s.click_cnt, 0)
		from banner b
			left join stat s on s.banner_id = b.id and s.slot_id = b.slot_id and s.usr_type_id = $2
		where b.slot_id = $1
		order by b.id
	`, pars.SlotID, pars.UsrTypeID)
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
			&item.ShowCnt,
			&item.ClickCnt,
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

func (d *St) BannerIncShowCount(ctx context.Context, pars entities.BannerStatIncPars) error {
	return d.bannerIncCol(ctx, "show_cnt", pars)
}

func (d *St) BannerIncClickCount(ctx context.Context, pars entities.BannerStatIncPars) error {
	return d.bannerIncCol(ctx, "click_cnt", pars)
}

func (d *St) bannerIncCol(ctx context.Context, col string, pars entities.BannerStatIncPars) error {
	v := pars.Value
	if v < 1 {
		v = 1
	}

	_, err := d.Db.ExecContext(ctx, `
		with q1 as (
		    update stat set `+col+` = `+col+` + $4
		    where banner_id = $1 and slot_id = $2 and usr_type_id = $3
		    returning banner_id
		)
		insert into stat(banner_id, slot_id, usr_type_id, `+col+`)
		select $1, $2, $3, $4
		where not exists(select * from q1)
	`, pars.ID, pars.SlotID, pars.UsrTypeID, v)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return err
	}

	return nil
}
