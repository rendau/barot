package pg

import (
	"context"
	"github.com/rendau/barot/internal/domain/entities"
)

func (d *St) StatIncShowCount(ctx context.Context, pars *entities.StatIncPars) error {
	return d.incCol(ctx, "show_cnt", pars)
}

func (d *St) StatIncClickCount(ctx context.Context, pars *entities.StatIncPars) error {
	return d.incCol(ctx, "click_cnt", pars)
}

func (d *St) incCol(ctx context.Context, col string, pars *entities.StatIncPars) error {
	_, err := d.db.ExecContext(ctx, `
		with q1 as (
		    update stat set `+col+` = `+col+` + 1
		    where banner_id = $1 and slot_id = $2 and usr_type_id = $3
		    returning banner_id
		)
		insert into stat(banner_id, slot_id, usr_type_id, `+col+`)
		select $1, $2, $3, 1
		where not exists(select * from q1)
	`, pars.BannerID, pars.SlotID, pars.UsrTypeID)
	if err != nil {
		d.lg.Errorw(ErrMsg, err)
		return err
	}

	return nil
}
