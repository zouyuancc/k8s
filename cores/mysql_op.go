package cores

import (
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Mysql_insert_task(data *Yaml) {
	// 连接数据库
	db, err := sqlx.Open("mysql", "test:123456@tcp(127.0.0.1:3306)/Kubernetes")
	if err != nil {
		fmt.Println("open mysql(Kubernetes) failed,", err)
		return
	} else {
		// 写入数据完成后关闭数据库
		defer func() {
			db.Close()
			fmt.Println("mysql(Kubernetes) closed successfully!")
		}()
		fmt.Println("mysql(Kubernetes) connected successfully!")
		// 创建唯一任务ID
		task_id := uuid.New().String()
		// 向task表中插入数据
		r, err := db.Exec("insert into task(task_id, kind, user, submit_time)values(?,?,?,?)", task_id, data.Kind, data.User, time.Now())
		if err != nil {
			fmt.Println("insert data into task table failed,", err)
			return
		}
		id, err := r.LastInsertId()
		if err != nil {
			fmt.Println("insert data into task table failed,", err)
			return
		}
		fmt.Println("task: insert succ:", id)
		if data.Kind == "Deployment" {
			// 向deployment表中插入数据
			metadata, _ := json.Marshal(data.Metadata)
			selector, _ := json.Marshal(data.Spec.Selector)
			template_metadata, _ := json.Marshal(data.Spec.Template.Metadata)
			r, err = db.Exec("insert into deployment(task_id, name, metadata, replicas, selector, template_metadata)values(?,?,?,?,?,?)", task_id, data.Metadata.Name, metadata, data.Spec.Replicas, selector, template_metadata)
			if err != nil {
				fmt.Println("insert data into deployment table failed,", err)
				return
			}
			id, err = r.LastInsertId()
			if err != nil {
				fmt.Println("insert data into deployment table failed,", err)
				return
			}
			fmt.Println("deployment: insert succ:", id)
			// 向container表中插入数据
			for _, value := range data.Spec.Template.Spec.Containers {
				container_id := uuid.New().String()
				command, _ := json.Marshal(value.Command)
				args, _ := json.Marshal(value.Args)
				resources, _ := json.Marshal(value.Resources)
				r, err = db.Exec("insert into container(container_id, task_id, image, name, command, args, resources)values(?,?,?,?,?,?,?)", container_id, task_id, value.Image, value.Name, command, args, resources)
				if err != nil {
					fmt.Println("container_id :", container_id, " insert data into container table failed,", err)
					return
				}
				id, err = r.LastInsertId()
				if err != nil {
					fmt.Println("container_id :", container_id, " insert data into container table failed,", err)
					return
				}
				fmt.Println("container_id ", container_id, " insert succ:", id)
			}
		}
		if data.Kind == "Service" {
			// 向service表中插入数据
			metadata, _ := json.Marshal(data.Metadata)
			selector, _ := json.Marshal(data.Spec.Selector)
			r, err = db.Exec("insert into service(task_id, name, metadata, type, selector)values(?,?,?,?,?)", task_id, data.Metadata.Name, metadata, data.Spec.Type, selector)
			if err != nil {
				fmt.Println("insert data into service table failed,", err)
				return
			}
			id, err = r.LastInsertId()
			if err != nil {
				fmt.Println("insert data into service table failed,", err)
				return
			}
			fmt.Println("service: insert succ:", id)
			// 向service_port表中插入数据
			for _, value := range data.Spec.Ports {
				p_id := uuid.New().String()
				r, err = db.Exec("insert into service_port(p_id, task_id, port, nodePort, targetPort)values(?,?,?,?,?)", p_id, task_id, value.Port, value.NodePort, value.TargetPort)
				if err != nil {
					fmt.Println("p_id ", p_id, " insert data into service_port table failed,", err)
					return
				}
				id, err = r.LastInsertId()
				if err != nil {
					fmt.Println("p_id ", p_id, " insert data into service_port table failed,", err)
					return
				}
				fmt.Println("service_port ", p_id, " insert succ:", id)
			}
		}
	}
}
